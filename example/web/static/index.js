(function () {
    this.title = "TinyAuth";
    this.version = "00.00.01";
    var signInDialog = document.getElementById('sign-in-dialog');
    var showSignInDialogBtn = document.getElementById("show-sign-in-dialog");
    showSignInDialogBtn.addEventListener('click', function () {
        signInDialog.showModal();
    });
    signInDialog.querySelector('.cancel').addEventListener('click', function () {
        signInDialog.close();
    });
    signInDialog.querySelector('.ok').addEventListener('click', function () {
        authenticate();
        signInDialog.close();
    });
    function onAuthenticateComplete(response) {
        ShowMessage("ok:" + (response.ok ? "Ok" : "Error") + ", status:" + response.status + ", status text:" + response.statusText, "Undo", function () { console.log("Undo"); });
    }
    function onAuthenticateError(err) {
        ShowMessage(err.message, "Undo", function () { console.log("Undo"); });
    }
    function authenticate() {
        var username = document.getElementById("username").value;
        var password = document.getElementById("password").value;
        var request = new Request('/authenticate', {
            method: 'GET',
            redirect: 'follow',
            headers: new Headers({
                'Authorization': "Basic " + btoa(username + ":" + password)
            })
        });
        fetch(request)
            .then(onAuthenticateComplete)
            .catch(onAuthenticateError);
    }
    var snackbarContainer = document.querySelector('#demo-snackbar-example');
    function ShowMessage(message, handlerLabel, handler) {
        var data = {
            message: message,
            timeout: 2000,
            actionHandler: handler,
            actionText: handlerLabel
        };
        snackbarContainer.MaterialSnackbar.showSnackbar(data);
    }
    rivets.bind(document.body, this);
}());
//# sourceMappingURL=index.js.map