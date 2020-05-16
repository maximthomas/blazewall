import React from 'react';
import { Button, TextField } from '@material-ui/core'

export class Callbacks extends React.Component {

    render() {
        
        const callbacks = this.props.callbacks.map(callback =>
            <div key={callback.name + "-container"}>
                <TextField error={!!callback.error} style={{width:"100%", marginTop: "14px"}} key={callback.name}
                    id={callback.name} type={callback.type} placeholder={callback.prompt}
                    onChange={(e) => this.props.updateCallback(e.target.value, callback.name)} helperText={callback.error}
                           value={callback.value} required={callback.required}/>
            </div>)
        const submitBtn = <div><Button  variant="contained" type="submit" color="default">Proceed</Button></div>
        const form = <form onSubmit={this.props.submitCallbacks} autoComplete={"off"}>{callbacks}{submitBtn}</form>
        return <div><h1>{this.props.title}</h1>{form}</div>
    }
}