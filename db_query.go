package main

import "fmt"

func WriteDB(data []byte) error {

	db, errCon := Connect()

	if errCon != nil {
		return fmt.Errorf("WriteDB__errCon %s", errCon)
	}

	defer db.Close()

	stmt, errQuery := db.Prepare(`insert into testing(reqjson) values(?)`)

	if errQuery != nil {
		return fmt.Errorf("WriteDB__errQuery %s", errQuery)
	}

	_, errData := stmt.Exec(data)

	if errData != nil {
		return fmt.Errorf("WriteDB__errData %s", errData)
	}

	return nil
}

func ReadDBGetID() (int, error) {
	db, errCon := Connect()

	if errCon != nil {
		return 0, fmt.Errorf("ReadDBGetID__errCon %s", errCon)
	}

	defer db.Close()

	var id int

	errQuery := db.QueryRow("select id from testing where id = (select max(id) from testing)").Scan(&id)

	if errQuery != nil {
		return 0, fmt.Errorf("ReadDBGetID__errQuery %s", errQuery)
	}

	return id, nil
}

func ReadDBGetResult(id string) ([]byte, error) {
	db, errCon := Connect()

	if errCon != nil {
		return nil, fmt.Errorf("ReadDBGetResult__errCon %s", errCon)
	}

	defer db.Close()

	var body []byte

	errQuery := db.QueryRow("select reqjson from testing where id =?", id).Scan(&body)

	if errQuery != nil {
		return nil, fmt.Errorf("ReadDBGetResult__errQuery %s", errQuery)
	}

	return body, nil
}
