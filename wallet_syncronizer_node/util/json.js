class Response {
    constructor(error, data) {
        this.error = error;
        this.data = data;
    }
}

function NewErrorResponse(httpStatus, message) {
    return new Response({
        status: httpStatus,
        message: message,
    }, null);
}

function NewSuccessResponse(data) {
    return new Response(null, data);
}

module.exports = {
    Response,
    NewErrorResponse,
    NewSuccessResponse,
};
