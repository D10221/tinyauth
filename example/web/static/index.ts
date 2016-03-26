/// <reference path="typings/tsd.d.ts"/>

interface SnackbarData {
    message: string;
    timeout: Number;
    actionHandler: ()=> void ,
    actionText: string
}

interface MaterialSnackbar {
    showSnackbar ( data: SnackbarData) ;
}

interface MaterialSnackBarContainer extends HTMLElement{
    MaterialSnackbar: MaterialSnackbar
}

(function () {

    this.title = "TinyAuth";

    this.version = `00.00.01`;

    var signInDialog = <HTMLElement>document.getElementById('sign-in-dialog');

    var showSignInDialogBtn = <HTMLButtonElement>document.getElementById("show-sign-in-dialog");

    showSignInDialogBtn.addEventListener('click', function() {
        (<any>signInDialog).showModal();
    });

    signInDialog.querySelector('.cancel').addEventListener('click', function() {
        (<any>signInDialog).close();
    });

    signInDialog.querySelector('.ok').addEventListener('click', function() {
        authenticate();
        (<any>signInDialog).close();
    });

    function onAuthenticateComplete(response: Response) {
        ShowMessage (
            `ok:${response.ok ? "Ok" : "Error"}, status:${response.status}, status text:${response.statusText}`,
            "Undo",
            ()=> { console.log("Undo")} );
        // responseOutput.innerHTML =
    }

    function onAuthenticateError(err: Error){
        ShowMessage (err.message, "Undo", ()=> { console.log("Undo")} )
    }

    function authenticate(/*event: Event*/){

        var username = (<HTMLInputElement>document.getElementById("username")).value;

        var password = (<HTMLInputElement>document.getElementById("password")).value;

        var request = new Request('/authenticate', {
            method: 'GET',
            // mode: 'cors',
            redirect: 'follow',
            headers: new Headers({
                //'Content-Type': 'text/plain',
                'Authorization': "Basic " + btoa(`${username}:${password}`)
            })
        });

        fetch(request)
            .then(onAuthenticateComplete)
            .catch(onAuthenticateError);
    }

    //... SnackBar
    var snackbarContainer = <MaterialSnackBarContainer>document.querySelector('#demo-snackbar-example');

    function ShowMessage(message: string, handlerLabel: string , handler : ()=> void ){
        var data = {
            message: message,
            timeout: 2000,
            actionHandler: handler,
            actionText: handlerLabel
        };
        snackbarContainer.MaterialSnackbar.showSnackbar(data);
    }

    // model-> UI : Binding
    rivets.bind(document.body, this );
}());


