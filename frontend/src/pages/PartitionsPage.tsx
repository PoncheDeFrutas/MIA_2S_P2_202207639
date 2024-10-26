import React from "react";
import { useParams } from "react-router-dom";
import Partition from "../components/Partition";

const partitionsData = [
    { id: "1", name: "Partition 1", fileSystem: "NTFS" },
    { id: "2", name: "Partition 2", fileSystem: "FAT32" },
    { id: "3", name: "Partition 3", fileSystem: "EXT4" },
];

const PartitionPage: React.FC = () => {
    const { diskId } = useParams<{ diskId: string }>();

    return (
        <div className="flex-grow flex items-center justify-center p-44">
            <div className="w-full max-w-3xl p-8 bg-white rounded-lg shadow-md">
                <h2 className="text-2xl font-bold mb-4 text-gray-800">
                    Partitions of Disk {diskId}
                </h2>
                <div className="flex flex-wrap gap-4">
                    {
                        partitionsData.map((partition) => (
                            <Partition
                                key={partition.id}
                                id={partition.id}
                                name={partition.name}
                            />
                        ))
                    }
                </div>
            </div>
        </div>
    );
}

export default PartitionPage;