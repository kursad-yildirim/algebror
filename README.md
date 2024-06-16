# algebror

Algebror is a simple GO program that generates primary school level math sheets for basic operations

## License
This project is licensed under the GNU General Public License v3.0. See the [LICENSE](./LICENSE) file for details.

## Properties
- Each sheet has 40 operations randomly selected over:
   - addition
   - subtraction
   - multiplication
   - division
   - percentage
- There are 5 difficulty levels as:
   - **Level-1:** Initial level, number range is short, decimals available only for addition and subtraction, mostly one significant digit after the decimal point
   - **Level-2:** Same as the first level with number range being a little more wider
   - **Level-3:** Number is wider than L2 and decimal numbers are available for multiplication with mostly one significant digit after the decimal point
   - **Level-4:** Number is wider than L3 and decimal numbers are available for division with mostly one significant digit after the decimal point
   - **Level-5:** Number range is wider than L4 and decimal numbers are available for percentage with mostly two significant digits after the decimal point
- It is either one of the operands or the result is asked from the student
- Application is starts a webserver on `:18080`
  - A `GET` request is accepted on path `/generate-test`
  - A pdf file of two pages is sent to the client. First page is the sent and the second page is the answer sheet.
  - Default difficulty level is `2`. if a parameter is sent within the url `?d=#` (where # must be an integer in [1:5]) to set the diffulty
  - A random index of five alphanumeric characters is added to each test
  - Pdf file is also stored in `out` directory within the application directory

## TODO
- Fractions test
- Exponential numbers
- Error handling
- No `out` directory
- Versioning
- Exception handling for `?d=#` when # > 5
- Different runtime environments
  - Docker
  - Kubernetes
  - Openshift
