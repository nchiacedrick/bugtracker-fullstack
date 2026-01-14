export interface Comment {
    id: string;
    bugId: string;
    author: string;
    content: string;
    createdAt: string;
}

export interface CreateCommentRequest {
    author: string;
    content: string;
} 