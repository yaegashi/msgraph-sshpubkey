import React from 'react';
import { Paper, Typography, Button, Theme } from '@material-ui/core';
import { useStyles } from './styles';
import { Link as RouterLink, LinkProps as RouterLinkProps } from 'react-router-dom';
import { makeStyles } from '@material-ui/styles';
import useRouter from 'use-react-router';
import { useTokenContext } from './token';

const Link1 = React.forwardRef<HTMLAnchorElement, RouterLinkProps>((props, ref) => <RouterLink innerRef={ref} {...props} />);

const useStyles2 = makeStyles((theme: Theme) => ({
    button: {
        marginRight: theme.spacing(1),
        textTransform: 'none',
    },
}))

export default () => {
    const c1 = useStyles();
    const c2 = useStyles2();
    const { getToken } = useTokenContext();
    const { location } = useRouter();
    const url = location.pathname && window.location.href;
    return (
        <Paper className={c1.paper}>
            <Typography variant="h6" color="inherit">
                Navigation
            </Typography>
            <p>
                <Button className={c2.button} variant="contained" color="default" component={Link1} to="/">/</Button>
                <Button className={c2.button} variant="contained" color="default" component={Link1} to="..">..</Button>
                <Button className={c2.button} variant="contained" color="default" component={Link1} to="A/">A</Button>
                <Button className={c2.button} variant="contained" color="default" component={Link1} to="B/">B</Button>
                <Button className={c2.button} variant="contained" color="default" component={Link1} to="C/">C</Button>
            </p>
            <p>
                <Button className={c2.button} variant="contained" color="primary" href={`/auth/sign_in?redirect=${url}`}>/auth/sign_in</Button>
                <Button className={c2.button} variant="contained" color="primary" href={`/auth/sign_out?redirect=${url}`}>/auth/sign_out</Button>
                <Button className={c2.button} variant="contained" color="primary" href="/auth/token">/auth/token</Button>
                <Button className={c2.button} variant="contained" color="primary" onClick={() => { getToken() }}>getToken()</Button>
            </p>
        </Paper>
    )
}