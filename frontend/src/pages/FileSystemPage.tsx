import File from "../components/File";
import Folder from "../components/Folder";
import { useParams } from "react-router-dom";
import React, { useState, useEffect } from "react";

const fileSystemData: { [key: string]: string } = {
    "1": "NTFS",
    "2": "FAT32",
    "3": "EXT4",
};

const FileSystemPage: React.FC = () => {
    const { partitionId } = useParams<{ partitionId: string }>();
    const [path, setPath] = useState("/");
    const [results, setResults] = useState<{ type: string; name: string }[]>([]);

    const fileSystem =
        partitionId && partitionId in fileSystemData
            ? fileSystemData[partitionId]
            : "Unknown";

    const handleSearch = (/*e: React.FormEvent<HTMLFormElement>*/) => {
        // TODO Make an API call to search for files and folders
        const simulateResults = [
            { type: "folder", name: "Documents" },
            { type: "folder", name: "Downloads" },
            { type: "file", name: "README.txt" },
            { type: "file", name: "index.txt" },
        ]
        setResults(simulateResults);
    }

    useEffect(() => {
        handleSearch();
    }, []);

    return (
        <div className="flex-grow flex flex-col items-center justify-center p-16">
            <div className="w-full max-w-3xl p-8 bg-white rounded-lg shadow-md">
                <h2 className="text-2xl font-bold mb-4 text-gray-800">
                    File System of Partition {partitionId}
                </h2>
                <p className="text-gray-700 mb-4">File System: {fileSystem}</p>
                <div className="flex mb-4">
                    <input
                        type="text"
                        value={path}
                        onChange={(e) => setPath(e.target.value)}
                        placeholder={`Search in ${path}`}
                        className="flex-grow p-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <button
                        onClick={handleSearch}
                        className="p-2 bg-blue-500 text-white rounded-r-md"
                    >
                        Search
                    </button>
                </div>
                <div className="flex flex-wrap gap-4">
                    {
                        results.map((result, index) => {
                            if (result.type === "folder") {
                                return <Folder key={index} name={result.name} />;
                            } else {
                                return <File key={index} name={result.name} />;
                            }
                        })
                    }
                </div>
            </div>
        </div>
    );
};

export default FileSystemPage;