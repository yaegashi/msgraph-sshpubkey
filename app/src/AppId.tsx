import React from 'react';
import { Paper, Typography } from '@material-ui/core';
import { useTokenContext } from './token';
import { useStyles } from './styles';
import SyntaxHighlighter from 'react-syntax-highlighter';

export default () => {
    const c1 = useStyles();
    const { token } = useTokenContext();
    const [id, setId] = React.useState(null);
    React.useEffect(() => {
        if (token) {
            fetch('https://graph.microsoft.com/v1.0/me', {
                method: 'GET',
                headers: new Headers({ Authorization: `Bearer ${token.access_token}` }),
            }).then(
                response => response.json()
            ).then(data => {
                setId(data);
            }).catch(error => {
                setId(null);
            })
        }
        else {
            setId(null);
        }
    }, [token]);
    return (
        <Paper className={c1.paper}>
            <Typography variant="h6" color="inherit">
                ID
            </Typography>
            {id &&
                <SyntaxHighlighter language="json">
                    {JSON.stringify(id, null, '  ')}
                </SyntaxHighlighter>
            }
        </Paper>
    )
}
