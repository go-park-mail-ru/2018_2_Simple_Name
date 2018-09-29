(function(){

class httpReq{
    _dofetch({url = '/',method = 'GET',data,callback = function(){}}={}){
        fetch(url, {
            method: method,
            mode: 'cors',
            credentials: 'include',
            headers: {
                "Content-Type": "application/json; charset=utf-8",
            },
            data: JSON.stringify(data),
        })
        .then(function(res){
            callback(res);
        });
    }
    doGet(params = {}){
        this._dofetch({...params,method: 'GET'});
    }
    
    doPost(params = {}){
        this._dofetch({...params,method : 'POST'});
    }
}
    window.httpModule = new httpReq();
})();