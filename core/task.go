package core

import (
	"database/sql"
	"math"
	"time"
)

type StartTaskRequest struct {
	QueueName string `validate:"required"`
	Data      []byte
	Time      time.Time `validate:"required"`
	Interval  int64     // Seconds
	Loop      int64
}

func StartTask(ctx *Context, request *StartTaskRequest) Error {
	ctx.LogInfo("Receive StartTaskRequest: %+v", *request)

	if err := validateStartTaskRequest(request); err != nil {
		ctx.LogError("Valudate task request fail: %v, err = %s", *request, err.Error())
		return ERROR_ADD_TASK_SYSTEM_FAIL
	}

	if request.Loop == -1 {
		request.Loop = math.MaxInt64
	}

	// Calculate new time
	nextTime := request.Time
	now := time.Now().Unix()
	requestTime := request.Time.Unix()
	var loopIndex int64 = 0
	if requestTime < now {
		loopIndex = (now - requestTime) / request.Interval
		if loopIndex > request.Loop {
			return ERROR_TASK_IS_EXPIRED
		}
		nextTimeSecond := (((now-requestTime)/request.Interval)+1)*request.Interval + requestTime
		nextTime = time.Unix(nextTimeSecond, 0)
	}
	bucket := GetBucket(nextTime) // Generate bucket id

	// Get id from id genrator
	taskId := ID.GenerateID()

	// Init transaction
	tx, err := DBSession().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		ctx.LogError("Begin transaction fail: %v, err = %s", *request, err.Error())
		return ERROR_ADD_TASK_SYSTEM_FAIL
	}
	defer tx.Rollback()

	// check task in database
	var id string
	var startTime string
	var loopCount, interval int64
	row := DBSession().QueryRowContext(ctx, "SELECT id, start_time, loop_count, interval FROM scheduler_tasks WHERE queue_name = $1", request.QueueName)
	err = row.Scan(&id, &startTime, &loopCount, &interval)
	if err == nil {
		if startTime == request.Time.Format(time.RFC3339) && loopCount == int64(request.Loop) && interval == request.Interval {
			return ERROR_TASK_ALREADY_EXISTED
		}

		ctx.LogInfo("Replace task: %s, startTime: %s, loopCount: %d, interval: %d", id, startTime, loopCount, interval)
		if _, err := tx.ExecContext(ctx, "DELETE FROM scheduler_tasks WHERE id = $1;", id); err != nil {
			ctx.LogError("delete task fail: %s, err = %s", id, err.Error())
			return ERROR_REMOVE_OLD_TASK_FAIL
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM scheduler_todo WHERE task_id = $1;", id); err != nil {
			ctx.LogError("delete todo task fail: %s, err = %s", id, err.Error())
			return ERROR_REMOVE_OLD_TASK_FAIL
		}
	}

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO scheduler_tasks(id, queue_name, data, done, loop_index, loop_count, next, interval, start_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
		taskId, request.QueueName, request.Data, false, loopIndex, request.Loop, nextTime.Unix(), request.Interval, request.Time.Format(time.RFC3339)); err != nil {
		ctx.LogError("Insert task fail: %v, err = %s", *request, err.Error())
		return ERROR_ADD_TASK_SYSTEM_FAIL
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO scheduler_todo(task_id, bucket) VALUES ($1, $2);", taskId, bucket); err != nil {
		ctx.LogError("Insert task fail: %v, err = %s", *request, err.Error())
		return ERROR_ADD_TASK_SYSTEM_FAIL
	}

	if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
		ctx.LogError("Commit fail: %v, err = %s", *request, err.Error())
		return ERROR_ADD_TASK_SYSTEM_FAIL
	}

	return nil
}

func validateStartTaskRequest(request *StartTaskRequest) Error {
	if err := validate.Struct(*request); err != nil {
		LoggerInstance.Error("Task request is invalid: %s", err.Error())
		return ERROR_TASK_REQUEST_INVALID
	}

	return nil
}

type StopTaskRequest struct {
	Id uint64
}

func StopTask(ctx *Context, request *StopTaskRequest) Error {
	tx, err := DBSession().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		ctx.LogError("Begin transaction fail: %v, error = %s", *request, err.Error())
		return ERROR_STOP_TASK_FAIL
	}

	// Delete todo in database
	if _, err := tx.ExecContext(ctx, "DELETE FROM todo WHERE id = $1", request.Id); err != nil {
		ctx.LogError("Delete task from todo fail: %d, %s", request.Id, err.Error())
		return ERROR_STOP_TASK_FAIL
	}

	// Delete task in database
	if _, err := tx.ExecContext(ctx, "DELETE FROM task WHERE id = $1", request.Id); err != nil {
		ctx.LogError("Delete task from tasks fail: %d, %s", request.Id, err.Error())
		return ERROR_STOP_TASK_FAIL
	}
	return nil
}
