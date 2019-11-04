import React from 'react';
import { Paper, Typography } from '@material-ui/core';
import { useTokenContext } from './token';
import { useStyles } from './styles';
import SyntaxHighlighter from 'react-syntax-highlighter';

export default () => {
    const c1 = useStyles();
    const { token } = useTokenContext();
    return (
        <Paper className={c1.paper}>
            <Typography variant="h6" color="inherit">
                Token
            </Typography>
            {token &&
                <SyntaxHighlighter language="json">
                    {JSON.stringify(token, null, '  ')}
                </SyntaxHighlighter>
            }
        </Paper>
    )
}
