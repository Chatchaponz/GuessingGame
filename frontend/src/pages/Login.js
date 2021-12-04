import Axios from "axios";
import React, { useState, useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../Component/Context";

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [loginStatus, setloginStatus] = useState('');
    const authToken = useContext(AuthContext)

    const loginInstance = Axios.create({
        baseURL: 'http://localhost:8080',
        headers: { "Login": "*" }
    });

    const navigate = useNavigate();

    // Set token in local storage
    const setLocalToken = (response) => {
        localStorage.setItem('token', response['data']);
        setloginStatus("")
        alert("Welcome " + username + " !!!")
        navigate("/");
    }

    const doLogin = async (event) => {
        event.preventDefault();
        await loginInstance.post('/login',
            JSON.stringify({
                username: username,
                password: password
            })
        ).then((response) => {
            console.log(response)
            setLocalToken(response)
        }).catch((error) => {
            console.error(error);
            setloginStatus("Invalid username or password!")
        });
    }

    useEffect(() => {
        if (authToken() === true)
            navigate("/");
    }, [authToken, navigate]);

    return (
        <div className="login_form">
            <form onSubmit={doLogin}>
                <div className="row justify-content-center">
                    <div className="col-md-3">
                        <br />

                        <h1 className="text-center">Login</h1>

                        <p className="text-center text-danger">{loginStatus}</p>

                        <input type="text" className="form-control" placeholder="Enter username" required
                            onChange={event => setUsername(event.target.value)}
                        />
                        
                        <br />

                        <input type="password" className="form-control" placeholder="Enter password" required
                            onChange={event => setPassword(event.target.value)}
                        />

                        <br />

                        <div className="text-center">
                            <button className="btn btn-primary" type="submit">Login</button>
                        </div>

                    </div>
                </div>
            </form>
        </div>
    );
};

export default Login;