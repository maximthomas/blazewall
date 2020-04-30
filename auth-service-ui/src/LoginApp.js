import React from 'react';
import { Callbacks } from './Callbacks'

const authUrl = process.env.REACT_APP_AUTH_URL

export class LoginApp extends React.Component {


    constructor(props) {
        super(props);
        this.state = {
            callbacks: [],
            module: null,
            succeeded: false,
            failed: false,
        };
        console.log(process.env);
        //this.updateCallback = this.updateCallback.bind(this);
    }
    componentDidUpdate(prevProps, prevState) {
        console.log(prevState, this.state);
    }

    componentDidMount() {
        this.getCallbacks()
    }


    getCallbacks() {
        fetch(authUrl, {
            credentials: "include",
        }).then((response) => {
            return response.json();
        }).then((data) => {
                if (data['callbacks']) {
                    this.setState({ callbacks: data['callbacks'] });
                }
                if (data['module']) {
                    this.setState({ module: data['module'] })
                }
            }).catch(function(e) {
                alert('Error connecting to a database')
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
        fetch(authUrl, {
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
                if (data['callbacks']) {
                    this.setState({ callbacks: data['callbacks'] });
                }
                if (data['module']) {
                    this.setState({ module: data['module'] })
                }
                if (data["status"]) {
                    if (data["status"] === "success") {
                        this.setState({ succeeded: true })
                    } else if (data["status"] === "failed") {
                        this.setState({ failed: true })
                    }
                }
            });

        return false;
    }

    render() {
        var uiComponent;
        if (this.state.succeeded) {
            uiComponent = <h1>Authentication succeeded</h1>
        } else if (this.state.failed) {
            uiComponent = <h1>Authentication failed</h1>
        } else {
            uiComponent = <Callbacks callbacks={this.state.callbacks}
                submitCallbacks={this.submitCallbacks}
                updateCallback={this.updateCallback} />
        }
        return <div id="login-app">
            {uiComponent}
        </div>
    };
}
