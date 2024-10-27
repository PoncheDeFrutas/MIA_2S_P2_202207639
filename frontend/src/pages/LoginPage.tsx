import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/authStore.ts";
import {POST} from "../services/API.ts";

const Login: React.FC = () => {
    const [partitionId, setPartitionId] = useState("");
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const navigate = useNavigate();
    const login = useAuthStore((state) => state.login);

    const handlePartitionIdChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPartitionId(e.target.value);
    };

    const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setUsername(e.target.value);
    };

    const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(e.target.value);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        try {
            const response = await POST<
                { partitionId: string, username: string, password: string },
                { result: boolean }
            >   ('login', { partitionId, username, password });
            if (!response.result) {
                throw new Error("Invalid credentials");
            }
            handleLogin(response.result);
        } catch (e) {
            alert(`Error: ${e}`);
        }
    };

    const handleLogin = (success: boolean) => {
        if (success) {
            login();
            navigate("/");
        } else {
            alert("Invalid credentials");
        }
    };

    return (
        <div className="flex justify-center items-center h-full p-32">
            <form
                onSubmit={handleSubmit}
                className="bg-white p-6 rounded shadow-md w-80"
            >
                <h2 className="text-2xl mb-4">Login</h2>
                <div className="mb-4">
                    <label className="block text-gray-700">ID Partition</label>
                    <input
                        type="text"
                        value={partitionId}
                        onChange={handlePartitionIdChange}
                        className="w-full px-3 py-2 border rounded"
                        required
                    />
                </div>
                <div className="mb-4">
                    <label className="block text-gray-700">Username</label>
                    <input
                        type="text"
                        value={username}
                        onChange={handleUsernameChange}
                        className="w-full px-3 py-2 border rounded"
                        required
                    />
                </div>
                <div className="mb-4">
                    <label className="block text-gray-700">Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={handlePasswordChange}
                        className="w-full px-3 py-2 border rounded"
                        required
                    />
                </div>
                <button
                    type="submit"
                    className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-700"
                >
                    Login
                </button>
            </form>
        </div>
    );
};

export default Login;