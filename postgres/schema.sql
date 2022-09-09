CREATE TABLE order_list
  (
     id                    SERIAL PRIMARY KEY,
     group varchar(50) NOT NULL,
     
     billing_id            VARCHAR(100) UNIQUE NOT NULL,
     seller_id             VARCHAR(100) NOT NULL,
     paid                  BOOLEAN DEFAULT false,
     tax_id                VARCHAR(8) NOT NULL DEFAULT '',
     area_code             VARCHAR(3) NOT NULL DEFAULT '',
     company               VARCHAR(40) NOT NULL DEFAULT '',
     address               VARCHAR(200) NOT NULL DEFAULT '',
     total                 MONEY NOT NULL,
     refund                MONEY NOT NULL,
     billing_type          INT2 NOT NULL DEFAULT 1,
     --{1: trial ,2: monthly ,3: annual ,4: commission ,5: monthly_commission ,6: annual_commission(commission) ,7; annual_commission(annual) }
     monthly_fee           MONEY NOT NULL,
     annual_fee            MONEY NOT NULL,
     commission_percentage FLOAT NOT NULL,
     create_at             DATE NOT NULL,
     transaction_id        VARCHAR(30),
     end_at                DATE,

     CONSTRAINT tax_id_valid CHECK (length(tax_id) = 8 OR length(tax_id) = 0),
     CONSTRAINT area_code_valid CHECK (length(area_code) = 3 OR length(area_code) = 0)
  );