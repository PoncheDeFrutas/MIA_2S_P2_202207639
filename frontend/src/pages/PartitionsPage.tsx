import {GET} from "../services/API.ts";
import { useParams } from "react-router-dom";
import React, {useEffect, useState} from "react";
import Partition from "../components/Partition";

type PartitionData = { id: string; name: string; fileSystem: string; };

const PartitionPage: React.FC = () => {
    const { diskId } = useParams<{ diskId: string }>();
    const [partitionsData, setPartitionsData] = useState<PartitionData[]>([]);

    useEffect(() => {
        (async () => {
            try {
                const response = await GET<{ result: PartitionData[] }>(`partitions/${diskId}`);
                if (!response.result) {
                    setPartitionsData([]);
                } else {
                    setPartitionsData(response.result);
                }

            } catch (e) {
                alert(`Error: ${e}`);
            }
        })();
    }, [diskId]);

    return (
        <div className="flex-grow flex items-center justify-center p-44">
            <div className="w-full max-w-3xl p-8 bg-white rounded-lg shadow-md">
                <h2 className="text-2xl font-bold mb-4 text-gray-800">
                    Partitions of Disk {diskId}
                </h2>
                <div className="flex flex-wrap gap-4">
                    {
                        partitionsData.length > 0 ? (
                        partitionsData.map((partition) => (
                            <Partition
                                key={partition.id}
                                id={partition.id}
                                name={partition.name}
                            />
                        ))
                        ) : (
                            <p className="text-gray-500">No partitions found</p>
                        )
                    }
                </div>
            </div>
        </div>
    );
}

export default PartitionPage;