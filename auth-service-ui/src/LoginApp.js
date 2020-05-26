import React from 'react';
import {Callbacks} from './Callbacks'
import {Profile} from './Profile'
import {ThemeProvider} from "@material-ui/styles";

import {createMuiTheme, CssBaseline, Link, Paper} from "@material-ui/core";

const signUpUrl = process.env.REACT_APP_SIGN_UP_URL
const signInUrl = process.env.REACT_APP_SIGN_IN_URL


const theme = createMuiTheme({
    palette: {
      type: "dark"
    }
  });

class AuthState {
    constructor(title, authUrl, linkTitle) {
        this.title = title;
        this.authUrl = authUrl;
        this.linkTitle = linkTitle;
    }
}

const signUpState = new AuthState("Sign Up", signUpUrl, "Sign In")
const signInState = new AuthState("Sign In", signInUrl, "Sign Up")

export class LoginApp extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            callbacks: [],
            module: null,
            succeeded: false,
            failed: false,
            authState: signInState,
        }
    }
    componentDidUpdate(prevProps, prevState, ss) {
        //console.log(prevState, this.state);
    }

    componentDidMount() {
        this.getCallbacks()
    }


    getCallbacks() {
        fetch(this.state.authState.authUrl, {
            credentials: "include",
        }).then((response) => {
                return response.json();
        }).then((data) => {
            this.processAuthData(data);
            }).catch((e) => {
                console.log(e)
                this.setState({ failed: true })
            });
        return []
    }

    updateCallback = (callbackValue, name) => {
        const callbacks = this.state.callbacks.slice();
        callbacks.forEach(callback => {
            if (callback.name === name) {
                callback.value = callbackValue;
            }
        });
        this.setState({ callbacks: callbacks });

    }

    submitCallbacks = (e) => {
        e.preventDefault();
        const callbacks = this.state.callbacks.slice();
        const request = {
            module: this.state.module,
            callbacks: callbacks,
        }
        const requestBody = JSON.stringify(request)
        fetch(this.state.authState.authUrl, {
            method: 'POST',
            body: requestBody,
            credentials: "include",
            headers: {
                'Content-Type': 'application/json'
            },
        }).then((response) => {
            return response.json();
        })
            .then((data) => {
                this.processAuthData(data);
            });

        return false;
    }

    switchAuthentication = () => {
        if(this.state.authState === signInState) {
            this.setState({authState: signUpState}, () =>this.getCallbacks());
        } else {
            this.setState({authState: signInState}, () => this.getCallbacks());
        }
    }

    processAuthData = (data) => {
        if (data['callbacks']) {
            this.setState({ callbacks: data['callbacks'] });
        }
        if (data['module']) {
            this.setState({ module: data['module'] });
        }
        if (data["status"]) {
            if (data["status"] === "success") {
                this.setState({ succeeded: true });
            }
            else if (data["status"] === "failed") {
                this.setState({ failed: true });
            }
        }
    }

    render() {
        let uiComponent;
        if (this.state.succeeded) {
            uiComponent = <div>
                <h1>Authentication succeeded</h1>
                <Profile/>
            </div>
        } else if (this.state.failed) {
            uiComponent = <h1>Authentication failed</h1>
        } else {
            uiComponent =<div><Callbacks callbacks={this.state.callbacks} title={this.state.authState.title}
                submitCallbacks={this.submitCallbacks}
                updateCallback={this.updateCallback} />
                <div id="links">
                    <Link id="auth-link" component="button" color="inherit" onClick={this.switchAuthentication}>
                        {this.state.authState.linkTitle}
                    </Link>
                </div>
            </div>
        }
        return <ThemeProvider theme={theme}>
            <CssBaseline/>
            <div id="login-app">
                <Paper id="auth-panel" variant="outlined">
                    {uiComponent}
                </Paper>
            </div>
        </ThemeProvider>;
    };
}
