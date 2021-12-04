import Axios from "axios";
import React, { useState, useContext, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../Component/Context";

const GuessGame = () => {
    const [guessNumber, setGuessNumber] = useState('');
    const [helpMessage, setHelpMessage] = useState('');
    const [infoMessage, setInfoMessage] = useState('');
    const [textColor, setTextColor] = useState('');

    const navigate = useNavigate();
    const authToken = useContext(AuthContext)
    const guessInstance = Axios.create({
        baseURL: 'http://localhost:8080',
        headers: { "Guess": "*" }
    });

    const logout = useCallback(() => {
        localStorage.clear();
        navigate("/Login");
    }, [navigate]);

    const handleGuessResponse = (response) => {
        let thisStatus = response['data']['status']
        setHelpMessage(`${guessNumber} is ${response['data']['result']}`);
        if (thisStatus === true) {
            setTextColor("text-success");
            setInfoMessage("New hidden number already random, you can guess again!");
        } else if (thisStatus === false) {
            setTextColor("text-danger");
            setInfoMessage("");
        } else {
            setTextColor("");
            setInfoMessage("");
        }
    }

    const doGuess = async (event) => {
        event.preventDefault();
        const token = localStorage.getItem('token');
        if (token == null) {
            alert("Token doesn't exist anymore")
            logout();
        }
        else {
            await guessInstance.post('/guess',
                JSON.stringify({
                    guess_number: parseInt(guessNumber)
                }),
                {
                    headers: { "Guess": "*", 'Authorization': `Bearer ${token}` }
                }
            ).then((response) => {
                console.log(response);
                handleGuessResponse(response);
            }).catch((error) => {
                console.error(error);
                alert(error)
                logout();
            });
        }
    }

    useEffect(() => {

        if (authToken() === true) {
            const thisToken = localStorage.getItem('token');
            Axios.get('http://localhost:8080/check_token', {
                headers: { 'Authorization': `Bearer ${thisToken}` }
            }).then((response) => {
                console.log(response)
            }).catch((error) => {
                console.error(error)
                alert(error)
                logout();
            });
        }
        else if (authToken() === false)
            logout();
        else
            console.log("Something worng with authToken in GuessGame")

    }, [authToken, logout]);

    return (
        <div className="guess_form">
            <form onSubmit={doGuess}>
                <div className="row justify-content-center">
                    <div className="col-md-4 text-center">
                        <br />

                        <h1>Guessing Game</h1>

                        <label htmlFor="GuessNumber" className="form-label">Guess the number! ( Between 0 - 100 )</label>

                        <p className={textColor}>{helpMessage}</p>

                        <input type="number" min="0" max="100" className="form-control" placeholder="Enter the number" required
                            onChange={event => setGuessNumber(event.target.value)}
                        />

                        <p className="text-success">{infoMessage}</p>

                        <button className="btn btn-primary" type="submit">Guess</button>

                        <hr />

                        <button className="btn" onClick={logout}>Logout</button>
                    </div>
                </div>
            </form>
        </div>
    );
};

export default GuessGame;