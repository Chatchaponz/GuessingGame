
import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import GuessGame from "./pages/GuessGame";
import { AuthContext } from "./Component/Context";

/* npm i react-router-dom */
/* npm i axios */
function App() {

    const authToken = () => {
        const thisToken = localStorage.getItem('token');
        if (thisToken == null)
            return false
        else {
            return true
        }
    }

    return (
        <div className="App container">
            <BrowserRouter>
                <AuthContext.Provider value={authToken}>
                    <Routes>
                        <Route path="/" exact element={<GuessGame />} />
                        <Route path="/Login" element={<Login />} />
                    </Routes>
                </AuthContext.Provider>
            </BrowserRouter>
        </div>
    );
}

export default App;
