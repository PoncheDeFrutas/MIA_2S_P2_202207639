export async function POST<T, R>(url: string, content: T): Promise<R> {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/`+url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(content),
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch ${url}`);
    }

    return await response.json() as Promise<R>;
}

export async function GET<R>(url: string): Promise<R> {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/`+url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch ${url}`);
    }

    return await response.json() as Promise<R>;
}
