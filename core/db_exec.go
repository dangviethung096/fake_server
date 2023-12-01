package core

import (
	"fmt"
	"reflect"
)

/*
* Save data to database
* @param data interface{} Data to save
* @return Error
 */
func SaveDataToDB[T DataBaseObject](ctx *Context, data T) Error {
	query, args, insertError := GetInsertQuery(data)
	if insertError != nil {
		ctx.LogError("Error when get insert data = %#v, err = %v", data, insertError)
		return insertError
	}

	ctx.LogInfo("Insert query = %v, args = %v", query, args)
	if _, err := sqliteSession.ExecContext(ctx, query, args...); err != nil {
		ctx.LogError("Error insert data = %#v, err = %v", data, err)
		return ERROR_SERVER_ERROR
	}

	return nil
}

/*
* Save data to database without primary key
* primary key will be auto increment in database
* @param data interface{} Data to save
* @return Error
 */
func SaveDataToDBWithoutPrimaryKey[T DataBaseObject](ctx *Context, data T) Error {
	query, args, pkAddress, insertError := GetInsertQueryWithoutPrimaryKey(data)
	if insertError != nil {
		ctx.LogError("Error when get insert data = %#v, err = %v", data, insertError)
		return insertError
	}

	ctx.LogInfo("Insert query = %v, args = %v", query, args)
	row := sqliteSession.QueryRowContext(ctx, query, args...)

	err := row.Scan(pkAddress)
	if err != nil {
		ctx.LogError("Get primary key from query fail: %v", err)
		return ERROR_INSERT_TO_DB_FAIL
	}

	return nil
}

/*
* Delete data in database
* @param data interface{} Data to delete
* @return Error
 */
func DeleteDataInDB[T DataBaseObject](ctx *Context, data T) Error {
	query, args, deleteError := GetDeleteQuery(data)
	if deleteError != nil {
		ctx.LogError("Error when get delete data = %#v, err = %v", data, deleteError)
		return deleteError
	}

	ctx.LogInfo("Delete query = %v, args = %v", query, args)
	if _, err := sqliteSession.ExecContext(ctx, query, args...); err != nil {
		ctx.LogError("Error delete data = %#v, err = %v", data, err)
		return ERROR_SERVER_ERROR
	}

	return nil
}

/*
* Update data in database
* @param data interface{} Data to update
* @return Error
 */
func UpdateDataInDB[T DataBaseObject](ctx *Context, data T) Error {
	query, args, updateError := GetUpdateQuery(data)
	if updateError != nil {
		ctx.LogError("Error when get update data = %#v, err = %v", data, updateError)
		return updateError
	}

	ctx.LogInfo("Update query = %v, args = %v", query, args)
	if _, err := sqliteSession.ExecContext(ctx, query, args...); err != nil {
		ctx.LogError("Error update data = %#v, err = %v", data, updateError)
		return ERROR_SERVER_ERROR
	}

	return nil
}

/*
* Select data from database by primary key
* @param data interface{} Data to select
* @return Error
 */
func SelectById(ctx *Context, data DataBaseObject) Error {
	query, params, err := GetSelectQuery(data)
	if err != nil {
		ctx.LogError("Error when get update data = %#v, err = %v", data, err)
		return err
	}

	query += fmt.Sprintf(" WHERE %s = $1", data.GetPrimaryKey())
	pk, found := searchPrimaryKey(data)
	if !found {
		ctx.LogError("Error not found primary key = %#v, err = %v", data, ERROR_NOT_FOUND_PRIMARY_KEY)
		return ERROR_NOT_FOUND_PRIMARY_KEY
	}

	ctx.LogInfo("Select query = %v, args = %v", query, pk.Interface())
	row := sqliteSession.QueryRowContext(ctx, query, pk.Interface())
	if err := row.Scan(params...); err != nil {
		ctx.LogError("Error select data = %#v, err = %v", data, ERROR_NOT_FOUND_IN_DB)
		return ERROR_NOT_FOUND_IN_DB
	}

	return nil
}

/*
* Select data from database by field: fieldName and fieldValue is passed in parameter
* @return Error
 */
func SelectByField(ctx *Context, data DataBaseObject, fieldName string, fieldValue any) Error {
	query, params, err := GetSelectQuery(data)
	if err != nil {
		ctx.LogError("Error when get update data = %#v, err = %v", data, err)
		return err
	}

	query += fmt.Sprintf(" WHERE %s = $1", fieldName)

	ctx.LogInfo("Select query = %v, args = %v", query, fieldValue)
	row := sqliteSession.QueryRowContext(ctx, query, fieldValue)
	if err := row.Scan(params...); err != nil {
		ctx.LogError("Error select data = %#v, err = %v", data, ERROR_NOT_FOUND_IN_DB)
		return ERROR_NOT_FOUND_IN_DB
	}

	return nil
}

/*
* searchPrimaryKey: search primary key in model
* @params: data DataBaseObject
* @return: reflect.Value, bool
 */
func searchPrimaryKey(data DataBaseObject) (reflect.Value, bool) {
	t, _ := getTypeOfPointer(data)
	found := false
	var idValue reflect.Value

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == data.GetPrimaryKey() {
			found = true
			idValue = reflect.ValueOf(data).Elem().FieldByIndex(field.Index)
			break
		}
	}
	return idValue, found
}
