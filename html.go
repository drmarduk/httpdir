package main

var (
	htmlHeader = `
		<html>
			<head>
				<title>Datenbunker</title>
			</head>
			<body>

	`

	htmlFooter = `
			</body>
		</html>
	`

	htmlLogin = htmlHeader + `
				<p>Login to the allmighty DATENBUNKER
					<form method="POST" action="/login">
						<input type="password" name="passphrase" />
						<input type="submit" value="GIMME GIMME GIMME" />
					</form>
				</p>
	` + htmlFooter
)
