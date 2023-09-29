### ESGI 4IBC | API CTF 2023 | GOLANG

---

# Port Scanning and Web Service Interaction Tool

This tool is designed to:

1. Scan a specific IP address for open ports within a specified range.
2. For each open port found, it sends HTTP requests to a list of known endpoints (paths) to fetch and display their responses.
3. Depending on the endpoint, the tool might POST data to the service, retrieve certain information (like a secret or user level), and use this information in subsequent requests to other endpoints.

## Components and Explanation:

### Data Structures:

- `Challenge`, `Content`, `SubmitSolutionBody`: Data structures representing the body of POST requests for specific endpoints.
- `ResponseDetails`: Struct to capture details of an HTTP response.
- `PostBody`: General struct used for most POST requests.

### Main Functions:

- `CheckPort(host string, port int, timeout time.Duration) bool`: Checks if a specific port is open on a host.
- `FindOpenPorts(host string, start, end int, timeout time.Duration) []int`: Uses goroutines to scan a range of ports concurrently and returns a list of open ports.
- `FetchDetails(ip string, port int, path string) ResponseDetails`: Makes a GET request to a given IP, port, and path, and returns the response details.
- `postBodyToCheckResponse(ip string, port int, path string, body interface{}) ([]byte, error)`: POSTs a JSON body to a specified endpoint and returns the response.

### Execution Flow:

1. Load configurations from a `.env` file: IP, start port, end port, and timeout.
2. Scan the IP for open ports within the given range.
3. For each open port found:
   - Fetch details from each known endpoint.
   - If the endpoint is `/getUserSecret`, a loop tries to fetch a secret until successful.
   - The secret fetched is then used for other endpoints.
   - If the endpoint is `/getUserLevel`, it fetches the user level.
   - If the endpoint is `/submitSolution`, it uses the previously fetched secret and user level to submit a solution.

## Usage:

1. Ensure you have a `.env` file in the root with the following variables: `IP`, `START_PORT`, `END_PORT`, and `TIMEOUT`.
2. Run `source .env` to load the variables into your environment.
3. Run the program using `go run main.go`.

## Note:

- The order of the paths in the list matters since some endpoints require information fetched from previous ones.
- The secret key for the `/submitSolution` endpoint currently has a humorous placeholder quote about pasting code. Ensure it's the right key for your application.

---

This README provides a high-level overview of your program. Depending on your audience and the purpose of this tool, you might want to expand on certain areas, add examples, or include other sections like "Prerequisites," "Installation," and "Contribution Guidelines."

## Challenge results:

### Indices

Tiny Path [ctf-school-????????] = ctf-school-09292023
Today is a good day innit ? = Today date

Copy Trash 5FPprcvF-T75f91DQ2C = url : https://pastebin.com/5FPprcvF , password : T75f91DQ2C

Dabatase App : 72 44 90 = Protocol
This is clearly not a binary : 81 49 56 53 50 51 53 = Q185235

### enterChallenge

Welcome to the challenge !
Here is your first Challenge:
77337396dc3250bc4c480e187a69b090
Don't forget that:Das Einfügen von Code aus dem Internet in Produktionscode ist ...

Secret Key = Pasting code from the Internet into production code is like chewing gum found in the street.