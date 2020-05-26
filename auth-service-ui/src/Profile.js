import React from 'react';

const idmUrl = process.env.REACT_APP_IDM_URL

export class Profile extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            session: {},
        }
    }

    componentDidMount() {
        this.getProfile();
    }

    getProfile = () => {
        fetch(idmUrl, {
            method: 'GET',
            credentials: "include",
            headers: {
                'Content-Type': 'application/json'
            },
        }).then((response) => {
            return response.json();
        }).then((data) => {
            this.setState({session: data})
        });
    }

    render() {
        if (!!this.state.session["Claims"]) {
            let props = this.state.session["Claims"]["props"];
            let propsElements = []
            if(!!props) {
                Object.keys(props).forEach((k) => {
                    propsElements.push(<p key={k}>{k}: {props[k]}</p>)
                });
            }
            return <div>
                <p>Login: {this.state.session["Claims"]["sub"]}</p>
                {propsElements}
            </div>
        }
        else  {
            return <div/>
        }
    }
}