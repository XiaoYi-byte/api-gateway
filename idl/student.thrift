namespace go demo

//--------------------request & response--------------
struct College {
    1: required string name(api.body = 'name'),
    2: string address(api.body = 'address'),
}

struct Student {
    1: required i32 id(api.body = 'id'),
    2: required string name(api.body = 'name'),
    3: required College college(api.body = 'college'),
    4: optional list<string> email(api.body = 'email'),
    5: required string sex(api.body = 'sex')
}

struct RegisterResp {
    1: bool success(api.body = 'success'),
    2: string message(api.body = 'message'),
}

struct QueryReq {
    1: required i32 id(api.body = 'id')
}

//----------------------service-------------------
service StudentService {
    RegisterResp Register(1: Student student)(api.post = '/student/register', api.param = 'true')
    Student Query(1: QueryReq req)(api.post = '/student/query', api.param = 'true')
}