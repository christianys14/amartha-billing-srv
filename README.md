# Billing Service
Dummy project related billing service for Amartha.

## Table of Content
- [Tech Stack](#tech-stack)
- [How to Run](#how-to-run)
- [How to Test](#how-to-test)
- [List of API](#list-of-apis)

## Tech Stack
1. Database MySQL
2. Golang 1.19
3. Cobra

## List of APIs
1. GET /v1/customer/outstanding/{customerID}
2. POST /v1/customer/payment
```json
{
   "user_id" : "f02b5a3f-692e-4c33-8ebd-5cc14afead73",
   "amount" : 4290000
}
```

## How to Run
I've 2 command which is :
1. "serveDummy" used for create data dummy insert into table.
2. "serveHttp" used for serve http rest api.

## How to Test
1. create database name with "amartha", and running migration scripts below : 
   - 20240623100359_create_table_loan.sql 
   - 20240623101036_create_index_loan.sql
   - 20240624043537_alter_table_loan.sql
2. update the detail of your database through credential.json (Mandatory)
3. update the detail of your port http through configuration.json (Optional)
4. my current IDE is using Intellij IDEA, if you're using also, click main.go and run. After that "edit configurations" from run menu, and see "program arguments" put the argument you would like to run.
5. run argument "serveDummy" first of all (Mandatory). Everytime this argument being triggered it would generated New of "customerID", etc.
6. run argument "serveHttp".
7. when the argument "serveDummy" being run, it would be insert as details below :
    - customer 1
      + w1 - w5 = paid
      + w6 - w50 = pending
    - customer 2
      + w1 - w50 = closed
    - customer 3
      + w1 - w50 = paid
8. Test cases :
    - Given the customers info below :
      + customer 1 = f02b5a3f-692e-4c33-8ebd-5cc14afead73
      + customer 2 = d480d9d7-ef39-4625-afaf-0363885688c1
      + customer 3 = 80e02637-b393-43fa-a8f0-d453d21faa36
      + customer 4 = c2ee4112-e00b-4d01-96ac-e3f83711ff2e
   - When invoke api get outstanding
```
//customer 1 will return :
{
    "rc": "0000",
    "message": "Successful",
    "data": {
        "remaining_outstanding": "4400000",
        "is_delinquent": true
    }
}

//customer 2,3 will return :
{
    "rc": "0000",
    "message": "Successful",
    "data": {
        "remaining_outstanding": "0",
        "is_delinquent": false
    }
}
```
   - When invoke api payment
```
//customer 1 
req
{
    "user_id" : "f02b5a3f-692e-4c33-8ebd-5cc14afead73",
    "amount" : 4290000
}

res
{
    "rc": "0000",
    "message": "Successful"
}
w6 - w45 status become PAID.

//customer 2, 3
req
{
    "user_id" : "d480d9d7-ef39-4625-afaf-0363885688c1",
    "amount" : 100000
}

res
{
    "rc": "0004",
    "message": "Congrats, you are not having any pending outstanding"
}

//customer 4 try to pay, but the amount is not equals with total outstanding (4290000)
req
{
    "user_id" : "c2ee4112-e00b-4d01-96ac-e3f83711ff2e",
    "amount" : 100000
}

res
{
    "rc": "0003",
    "message": "amount of payment should be exact"
}
```