package cmds

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nwillems/glauth-tools/chpass"
)

// PasswordChangeServer start server that will allow changing password
func PasswordChangeServer(config chpass.Config, args []string) {
	http.HandleFunc("/", serveWebPage())
	http.HandleFunc("/password", changePassword(config))

	log.Print("Starting HTTP server on :9003")
	log.Fatal(http.ListenAndServe(":9003", nil))
}

func serveWebPage() http.HandlerFunc {
	fmt.Print("Serving the static page")
	const passwordPage = `<!DOCTYPE html>
	<html lang="en">
	<html>
	
	<head>
		<link rel="stylesheet" href="https://unpkg.com/mvp.css">
	</head>
	
	<body>
		<main>
			<hr />
			<section>
				<form action="/password" method="post">
					<header>
						<h2>Change Password</h2>
					</header>
					<label for="username">Username</label>
					<input type="text" name="username" />
					<label for="old_password">Old password</label>
					<input type="password" name="old_password" />
	
					<label for="new_password">New password</label>
					<input type="password" name="new_password" />
					<label for="new_password_repeat">Repeat new password</label>
					<input type="password" name="new_password_repeat" />
	
					<button type="submit">Change password</button>
				</form>
			</section>
		</main>
	</body>
	
	</html>
`

	return func(rw http.ResponseWriter, r *http.Request) {
		// http.ServeFile(rw, r, "password.html")
		rw.Header().Add("Content-Type", "text/html")
		rw.WriteHeader(200)
		fmt.Fprint(rw, passwordPage)
	}
}

func changePassword(config chpass.Config) http.HandlerFunc {
	fmt.Print("Changin password")

	return func(rw http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")

		rw.Header().Add("Content-Type", "text/plain")

		if newPassword != r.FormValue("new_password_repeat") {
			rw.WriteHeader(403)
			fmt.Fprintln(rw, "New password and repeated password were not equal")
			return
		}

		err := chpass.ChangePassword(username, oldPassword, newPassword, config)
		if err != nil {
			rw.WriteHeader(500)
			fmt.Fprintf(rw, "Password not changed - %s\n", err.Error())
		}

		rw.WriteHeader(200)
		fmt.Fprintln(rw, "Password changed!")
	}
}
