# Go Wallet

## 1. How it should work?

- Monzo fires a webhook when a ransaction is made
- The API get's the transaction details and creates a record in the db
- A message is sent to the telegram app asking what did you buy
- User inputs his list of stuff that he bought
- DB record is updated with the details
- If user does not reply to the message, another telegram message will be sent every hour until the end of day
- Report of transactions is generated per day / week / month

## 2. What does the report contain?

- What did you buy in this period of time
- How does this compare to the previous period of time
- How much money did you save or went over

## 3. Database schema

- transactions
    - id
    - monzo_id
    - merchant_id
    - amount
    - currency
    - description
    - notes
    - category

- accounts
    - id
    - monzo_id
    - balance
    - currency

- receipts
    - id
    - transaction_id
    - description
    - amount
    - quantity
    - unit