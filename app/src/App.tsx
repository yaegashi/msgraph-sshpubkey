import React from 'react';
import './App.css';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import { BrowserRouter } from 'react-router-dom';
import { TokenProvider, useTokenContext } from './token';
import AppBar from './AppBar';
import AppNav from './AppNav';
import AppToken from './AppToken';
import AppId from './AppId';

const theme = createMuiTheme();

const AppGetToken = () => {
    const { token, getToken } = useTokenContext();
    if (token == null) {
        getToken();
    }
    return null;
}

export default () => {
    return (
        <TokenProvider>
            <AppGetToken />
            <MuiThemeProvider theme={theme}>
                <BrowserRouter basename={process.env.PUBLIC_URL}>
                    <AppBar />
                    <AppNav />
                    <AppToken />
                    <AppId />
                </BrowserRouter>
            </MuiThemeProvider>
        </TokenProvider>
    );
}
