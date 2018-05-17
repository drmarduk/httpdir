package main

var (
	htmlHeader = `<!DOCTYPE html>
		<html>
			<head>
				<title>Datenbunker</title>
				<style>
					body, html {
						font-family: helvetica;
					}

					li {
						list-style: none;
						height: 15px;
						width: 900px;
						/*background: #DDDDDD;*/
						padding: 7px;
						margin-top: 3px;
						margin-right: 50px;
					}
					li:hover {
						background: #f5f5f5;
					}
				</style>
			</head>
			<body>
	`

	htmlFooter = `
			</body>
		</html>
	`

	htmlLogin = htmlHeader + `
				<p id="login_a">Login to the allmighty DATENBUNKER
					<form method="POST" action="/login">
						<input type="password" name="passphrase" />
						<input type="submit" value="GIMME GIMME GIMME" />
					</form>
				</p>
	` + htmlFooter
)
