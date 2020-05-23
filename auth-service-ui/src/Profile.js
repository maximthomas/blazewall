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
            return <div>
                <p>Login: {this.state.session["Claims"]["sub"]}</p>
                <p>Name: {this.state.session["Claims"].props?.name}</p>
            </div>
        }
        else  {
            return <div/>
        }
    }
}