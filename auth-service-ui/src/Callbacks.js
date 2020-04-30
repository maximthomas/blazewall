import React from 'react';
export class Callbacks extends React.Component {

    render() {
        
        const callbacks = this.props.callbacks.map(callback =>
            <div>
                <input key={callback.name} id={callback.name}
                    type={callback.type} placeholder={callback.prompt}
                    onChange={(e) => this.props.updateCallback(e.target.value, callback.name)} />
                <span>{callback.error}</span>
            </div>)
        const submitBtn = <input type="button" value="Submit" onClick={this.props.submitCallbacks}></input>

        const form = <form onSubmit={this.props.submitCallbacks}><input type="text"></input>{callbacks}{submitBtn}</form>
        return <div><h1>Authenticate</h1>{form}</div>
    }
}