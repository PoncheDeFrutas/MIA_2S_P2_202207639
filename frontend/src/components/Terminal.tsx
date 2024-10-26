import { EditorComponent } from "./Editor";
import React, { useState } from 'react';
import {POST} from "../services/API.ts";

export default function Terminal() {
    const [inputText, setInputText] = useState('');
    const [outputText, setOutputText] = useState('');
    const [fileName, setFileName] = useState('');

    const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            setFileName(file.name);
            reader.onload = (e) => {
                setInputText(e.target?.result as string);
            };
            reader.readAsText(file);
        }
    };

    const handleSaveFile = () => {
        const blob = new Blob([outputText], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'output.smia';
        a.click();
        URL.revokeObjectURL(url);
    };

    const handleExecute = async () => {
        try {
            const response = await POST<
                { content: string },
                { Result: string }
            >   ('http://localhost:5000/execute', { content: inputText }
            );
            setOutputText(response.Result);
        } catch (e) {
            setOutputText(`Error: ${e}`);
        }
    };

    const handleContentChange = (content: string) => {
        setInputText(content);
    };

    return (
        <div className="App h-screen">
            <header className="flex justify-between items-center p-4 bg-gray-800 text-white">
                <div className="header-left">
                    <h2 className="text-xl">Terminal</h2>
                </div>
                <div className="header-right flex items-center">
                    <div className="header-info mr-4">
                        <h2 className="text-lg">{fileName || "No file selected"}</h2>
                    </div>
                    <div className="header-buttons flex space-x-2">
                        <div className="file-upload-wrapper relative">
                            <input
                                type="file"
                                id="file-upload"
                                accept=".smia"
                                onChange={handleFileUpload}
                                className="file-input absolute inset-0 opacity-0 cursor-pointer"
                            />
                            <button
                                className="header-button bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                                onClick={() => document.getElementById('file-upload')?.click()}
                            >
                                Upload File
                            </button>
                        </div>
                        <button
                            className="header-button bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
                            onClick={handleExecute}
                        >
                            Execute
                        </button>
                        <button
                            className="header-button bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600"
                            onClick={handleSaveFile}
                        >
                            Save
                        </button>
                    </div>
                </div>
            </header>
            <div className="content flex flex-row  h-[75%]">
                <div className="left-side w-1/2 p-4">
                    <EditorComponent
                        content={inputText}
                        setContent={handleContentChange}
                    />
                </div>
                <div className="right-side w-1/2 p-4">
                <textarea
                    className="output-area w-full h-full border border-gray-300 rounded p-2"
                    readOnly
                    value={outputText}
                    placeholder="Output will appear here..."
                ></textarea>
                </div>
            </div>
        </div>
    );

}
