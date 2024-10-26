import {BrowserRouter, Route, Routes} from "react-router-dom";

import TerminalPage from "./pages/TerminalPage";
import LoginPage from "./pages/LoginPage.tsx";
import DiskPage from "./pages/DiskPage.tsx";
import PartitionsPage from "./pages/PartitionsPage.tsx";
import FileSystemPage from "./pages/FileSystemPage.tsx";

import NavBar from "./components/NavBar";
import Footer from "./components/Footer.tsx";

export default function App() {
    return (
        <BrowserRouter>
            <div className="flex flex-col min-h-screen bg-gray-100">
                <NavBar/>
                <main className="container mx-auto px-4 flex-grow">
                    <Routes>
                        <Route path="/" element={<TerminalPage/>}/>
                        <Route path="/login" element={<LoginPage/>}/>
                        <Route path="/disks" element={<DiskPage/>}/>
                        <Route path="/partitions/:diskId" element={<PartitionsPage/>}/>
                        <Route path="/filesystem/:partitionId" element={<FileSystemPage/>}/>
                    </Routes>
                </main>
                <Footer/>
            </div>
        </BrowserRouter>
    );
}