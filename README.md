Authentication

Register
Create your account by providing your name, email, password, phone number, and role. Once registered, you’ll receive an authentication token that you’ll use for future requests.

Login
Enter your email and password to access your account. If the credentials are correct, a new authentication token will be generated for you.

Logout
Logs you out of the system. After logging out, your current token will no longer be valid.

Employee Management

Create Employee (Admin Only)
Admins can create employee accounts directly instead of requiring employees to self-register.

List Employees with Asset Count
View all employees along with the number of assets currently assigned to them. You can filter results by asset type (e.g., laptop) or asset status (e.g., assigned). If no filters are applied, it displays the total assets per employee.

Archive Employee (Admin Only)
Marks an employee as inactive (for example, if they leave the company). The employee is not permanently deleted — their data remains in the system.

My Assets
Allows an employee to see all assets currently assigned to them. The system identifies the employee automatically using their authentication token.

Asset Management

Create Asset (Admin Only)
Adds a new device to the system. Basic details like brand, model, serial number, and type are required. Depending on the asset type, additional specifications must be provided — for example, RAM and storage for laptops, or DPI for a mouse.

List All Assets
Displays all company assets. Supports searching, filtering (by type, status, or owner), and pagination. It also provides dashboard statistics such as how many assets are available, assigned, damaged, etc.

Get One Asset
Retrieves detailed information about a specific asset, including type-specific specifications like processor, RAM, and storage for laptops.

Assign Asset (Admin Only)
Assigns an asset to an employee. The asset status changes to assigned, and the action is recorded in the asset history.

Return Asset (Admin Only)
Marks an asset as returned. Its status changes back to available, and a return record is saved in the history along with an optional note.

Update Asset (Admin Only)
Allows modification of asset details such as brand, model, warranty dates, or technical specifications.

Delete Asset (Admin Only)
Soft-deletes an asset from the system. The asset must not be currently assigned. The data remains in the database but is hidden from active listings.
