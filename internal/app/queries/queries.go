package queries

const Register = `INSERT INTO Usr(User, Mail, Tel, Pwd, Date_reg, St) VALUES (?, ?, ?, ?, NOW(), 'ACTIVE');`

const Search = `SELECT id_Usr, User, Mail, Tel, Pwd, Date_reg, St FROM Usr Where (St='ACTIVE') AND (Mail = ? OR Tel = ?) `
const ExistUsr = `SELECT id_Usr, St FROM Usr Where (St='ACTIVE') AND (Mail = ? OR Tel = ?) `

const Cancel = `UPDATE Usr SET St='CANCEL' WHERE id_Usr = ?`
