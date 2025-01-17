# Granitex
_italic A secure message transmission system_
Developed by: JD Sacharok [https://github.com/Gorillarock](https://github.com/Gorillarock)

## Dependecies
1. docker
- For installation instructions follow: [https://docs.docker.com/get-started/get-docker/](https://docs.docker.com/get-started/get-docker/)
2. make
- For Linux:
    ```
    sudo apt-get install make
    ```
    _italic NOTE: may use different package manager, such as apt, if your distro is different._
- For Mac:
    ```
    brew install make
    ```
    _italic NOTE: will require brew package manager  [https://brew.sh/](https://brew.sh/)_
- For Windows:
    Follow instructions here [https://gnuwin32.sourceforge.net/packages/make.htm](https://gnuwin32.sourceforge.net/packages/make.htm)


## Getting Started
1. Create a .env file in the same directory as the existing .env.dist file.
2. Provide values for all of the ENV keys in the .env file.
    _italic NOTE: Choose what you want, because the MongoDB and Server will each use the provided values._
3. Run Command:
    ```
    make run
    ```
4. In a browser, navigate to [http://localhost:80/v1/tx](http://localhost:80/v1/tx)

## Security Concepts Used
### Clients
- Clients will *NEVER* transmit the PIN to the server.
    - Preventing the server from being able to decrypt the message.
- Clients will only ever recieve the encrypted message from the server when ID, VERIFY, and ANSWER all result in a match.
    - This effectively prevents any chance at unauthorized reception of the encrypted message.
    - Since the encrypted message is only returned to the authorized accessor, there is no way to brute force attack the encrypted message.

### Server
- *NEVER* has access to the PIN.
    - Preventing the server from being able to decrypt the message.
- *NEVER* stores the QUESTION.
    - Preventing brute force attacks against entries in the DB in the case of data breach.
    - The brute force being prevented involves iteratively appending entries from a rainbow table, as the PIN, (dictionary attack) to the question and hashing to try to obtain the correct answer, and then using the entry as the PIN to decrypt the ill-begotten encrypted message.
- Uses authentication to access the database to reduce chances of data breach.
    _italic NOTE: This is basic auth which is not secure unless using TLS to communicate to the server (not advised for production usage)_
- Uses and ID (UUID) and VERIFY code match to validate that the DB entry should be checked (including incrementation of tries).
    - Rendering denial of service (DOS) type attacks, designed to wipe DB entries through repeated random attempts all but completely useless.
    - Not only are UUID matches unlikely, but by matching with a cryptographically secure VERIFY code (salt), the likelyhood of guessing both an existing UUID and the matching random VERIFY (which is a number randomly generated between 1 and 1,000,000) is astronomically small.
    - Chance of guessing an ID grows as the DB entries grow, but guessing any in particular is a _italic 1 in 3.4×10^38_.
    - Assuming an entry ID is guessed, the chance of matching the VERIFY adds another _italic 1 in 1 x 10^6_ at further guessing the matching VERIFY.
    - Since status 404 is returned for any incorrect ID and VERIFY combination which does not exist in the DB, attackers can only even attempt the PIN if they guessed exactly the ID and the VERIFY.
- After 3 unauthorized attempts at the ANSWER to a given DB entry match, the entrie entry is deleted, and any further attempts will return status 404.
- After successful authorized access to an entry, the server returns only the encrypted message, and then the server deletes the entry, and any further attempts will return status 404.

### DB
- ENV vars for DB credentials means that the DB credentials and not stored in github.


## Clients
### v1/tx client
#### Process
![TX client / Server Communication](./diagrams/c1_serv_swim.jpg)

#### Usage
NOTE: At any time, the granitex logo can be clicked to refresh the Tx client.
1. Fill out fields:
    - Fields
        - Message: Enter a message which you would like to transmit in secret.
        - PIN:     A secret code to be used for encrypting/decrypting the encrypted message.
2. Press _italic Submit_ button
NOTE: Server will never store the question.
NOTE 2: TX Client will *NEVER* transmit the PIN to the server.

### v2/rx client
#### Process
![RX client / Server Communication](./diagrams/c2_serv_swim.jpg)

#### Usage
NOTE: At any time, the granitex logo can be clicked to navigate to the Tx client.
1. Fill out fields:
    - Fields:
        - PIN: The secret code to be used for decrypting the encrypted message.
        _italic NOTE: The same one used when submitting with tx client._
2. Press _italic Get_ button
NOTE: The RX Client will *NEVER* transmit the PIN to the server.

## Deployment Notes
_italic Advisements before deploying for production usage._
- Add TLS options to Server for Production Deployment.
- Use a dedicated Atlas or self-hosted MongoDB solution (docker-compose is only advised for local use).


## Development
### Nice to haves
- A user management system:
    - Log in to manage stored messages.
    - Enables billing based on usage.
- Time to live for all entries (unless logged into proper tier).
- Tooltips to for fields which explain their usage.
