import React from 'react';
import Editor from '@monaco-editor/react';

interface EditorComponentProps {
    content: string;
    setContent: (content: string) => void;
}

export class EditorComponent extends React.Component<EditorComponentProps> {
    handleEditorChange = (value: string | undefined) => {
        this.props.setContent(value || '');
    }

    render() {
        const { content } = this.props;
        return (
            <div className="w-full h-full flex flex-col">
                <div className="w-full h-full">
                    <Editor
                        theme="vs-dark"
                        language="Markdown"
                        value={content}
                        options={{
                            wordWrap: 'on',
                            minimap: { enabled: false },
                            scrollBeyondLastLine: false,
                            automaticLayout: true,
                        }}
                        onChange={this.handleEditorChange}
                    />
                </div>
            </div>
        );
    }
}