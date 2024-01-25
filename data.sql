CREATE TABLE Wallet (
                        "id" uuid primary key default gen_random_uuid(),
                        "balance" int not null
);

CREATE TABLE TransferHistory(
                                "id" uuid primary key default gen_random_uuid(),
                                "time" timestamp not null,
                                "fromwallet" text not null ,
                                "towallet" text not null ,
                                "amount" int not null
);

INSERT INTO wallet (balance) VALUES (100)

SELECT id, balance  FROM wallet

SELECT id,fromwallet,towallet,amount FROM transferhistory