(function(){

class httpReq{
    _dofetch({url = '/',method = 'GET',data,callback = function(){}}={}){
        fetch(url, {
            method: method,
            mode: 'cors',
            credentials: 'include',
            //headers: {
            //    "Content-Type": "multipart/form-data; boundary=something",//"application/json; charset=utf-8",
          //  },
            body: data,
        })
        .then(function(res){
            callback(res);
        })
        .catch(function(err){
            console.log(err);
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
