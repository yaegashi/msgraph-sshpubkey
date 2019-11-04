import React from 'react';

type TokenContextProps = { token: null | { access_token: string }, getToken: () => void };

const initialContext: TokenContextProps = { token: null, getToken: () => { } };

export const TokenContext = React.createContext(initialContext);

export const TokenProvider = (props: { children: React.ReactNode; }) => {
    const [token, setToken] = React.useState(null);
    const getToken = () => {
        fetch('/auth/token').then(
            response => response.json()
        ).then(data => {
            setToken(data);
        }).catch(error => {
            setToken(null);
        });
    }
    return (
        <TokenContext.Provider value={{ token, getToken }}>
            {props.children}
        </TokenContext.Provider>
    )
}

export const useTokenContext = () => React.useContext(TokenContext);
