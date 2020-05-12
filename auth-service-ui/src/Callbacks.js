import React from 'react';
import { Button, TextField } from '@material-ui/core'

export class Callbacks extends React.Component {

    render() {
        
        const callbacks = this.props.callbacks.map(callback =>
            <div>
                <TextField error={!!callback.error} style={{width:"100%", marginTop: "14px"}} key={callback.name} id={callback.name}
                    type={callback.type} placeholder={callback.prompt}
                    onChange={(e) => this.props.updateCallback(e.target.value, callback.name)} 
                    helperText={callback.error}/>
            </div>)
        const submitBtn = <div><Button  variant="contained" type="submit" color="default">Login</Button></div>
        const form = <form onSubmit={this.props.submitCallbacks}>{callbacks}{submitBtn}</form>
        return <div><h1>Authentication</h1>{form}</div>
    }
}