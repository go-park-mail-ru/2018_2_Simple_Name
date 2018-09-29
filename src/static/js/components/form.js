class baseView {
    render(){};
    show(){};
    hide(){};
}

export class Form{
    constructor(){
        this.prototype = Object.create(baseView.prototype);
    }
    render({ inputs = {}, formId = '' }) {
        const form = document.createElement('form');
        form.id = formId;
        inputs.forEach(function (item) {
            const input = document.createElement('input');
    
            input.name = item.name;
            input.type = item.type;
            if (item.value) {
                input.value = item.value;
            }
            input.placeholder = item.placeholder;
    
            form.appendChild(input);
            form.appendChild(document.createElement('br'));
        });
        return form;
    }
}