import Disk from "../components/Disk";

const disksData = [
    { id: "1", name: "Disk 1" },
    { id: "2", name: "Disk 2" },
    { id: "3", name: "Disk 3" },
];

function FileExplorer() {
    return (
        <div className="flex-grow flex items-center justify-center p-44">
            <div className="w-full max-w-3xl p-8 bg-white rounded-lg shadow-md">
                <h2 className="text-2xl font-bold mb-4 text-gray-800">Disks</h2>
                <div className="flex flex-wrap gap-4">
                    {
                        disksData.map((disk) => (
                            <Disk key={disk.id} id={disk.id} name={disk.name}/>
                        ))
                    }
                </div>
            </div>
        </div>
    );
}

export default FileExplorer;
