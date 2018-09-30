'use strict';

import { Form } from "./components/form.js";

const httpReq = window.httpModule;

const root = document.getElementById("root");

function createLinkMenu() {
    const link = document.createElement('a');
    link.textContent = "Back to main menu";
    link.href = "menu";
    link.dataset.href = "menu";
    return link;
}

function createMenu() {
    const logo = document.createElement('div');
    logo.classList.add("p_name");

    const logo_header = document.createElement('h1');
    logo_header.innerHTML = 'Simple game';
    logo.appendChild(logo_header);

    const titles = {
        signin: 'Sign in',
        signup: 'Sign up',
        leaders: 'Leaders',
        profile: 'Profile',
        about: 'About'
    };

    const dl = document.createElement('dl');
    dl.classList.add("menu");

    Object.entries(titles).forEach(function (entry) {
        const dt = document.createElement('dt');
        dt.classList.add("button_menu");

        const href = entry[0];
        const title = entry[1];

        const a = document.createElement('a');
        a.href = href;
        a.dataset.href = href;
        a.title = title;
        a.textContent = title;

        a.classList.add("button_menu");

        dt.appendChild(a);
        dl.appendChild(dt);



    });
    root.appendChild(logo);

    root.appendChild(dl);


}

function createSignIn() {
    root.innerHTML = '';

    const header = document.createElement('div');
    header.dataset.sectionName = 'header';
    header.id = "header";
    root.appendChild(header);

    const logo = document.createElement('span');
    logo.id = 'logo';
    const headerTitle = document.createElement('h1');
    headerTitle.textContent = 'Sign in';
    headerTitle.id = "headerTitle";

    header.appendChild(logo);
    header.appendChild(headerTitle);


    const body = document.createElement('div');
    body.id = 'body';
    root.appendChild(body);

    const formblock = document.createElement('div');
    formblock.id = 'formblock';
    body.appendChild(formblock)

    const inputs = [
        {
            name: 'email',
            type: 'email',
            placeholder: 'Email'
        },
        {
            name: 'password',
            type: 'password',
            placeholder: 'Password'
        },
        {
            name: 'submit',
            type: 'submit',
            value: 'Submit'
        }
    ];
    const formObj = new Form();
    const form = formObj.render({ inputs: inputs, formId: 'signinForm' })
    formblock.appendChild(form);

    const pLink = document.createElement('p');
    form.appendChild(pLink);
    const link = createLinkMenu()
    pLink.appendChild(link);

    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const email = form.elements['email'].value;
        const password = form.elements['password'].value;
        if (email == "" || password == "") {
            alert("Enter email!")
            return
        }
        httpReq.doPost({
            url: '/signin',
            callback(res) {
                if (res.status == 404) {
                    alert("Wrong login or password");
                    return;
                }
                if(res.status == 400){
                    alert("Already loged in")
                }
                createProfile();
            },
            data: {
                email,
                password
            }
        })
    });
}

function createSignUp() {
    root.innerHTML = '';

    const header = document.createElement('div');
    header.dataset.sectionName = 'header';
    header.id = "header";
    root.appendChild(header);

    const logo = document.createElement('span');
    logo.id = 'logo';
    const headerTitle = document.createElement('h1');
    headerTitle.textContent = 'Sign up';
    headerTitle.id = "headerTitle";

    header.appendChild(logo);
    header.appendChild(headerTitle);


    const body = document.createElement('div');
    body.id = 'body';
    root.appendChild(body);

    const formblock = document.createElement('div');
    formblock.id = 'formblock';
    body.appendChild(formblock)
    const inputs = [
        {
            name: 'email',
            type: 'email',
            placeholder: 'Email'
        },
        {
            name: 'password',
            type: 'password',
            placeholder: 'Password'
        },
        {
            name: 'password_repeat',
            type: 'password',
            placeholder: 'Repeat Password'
        },
        {
            name: 'submit',
            type: 'submit',
            value: 'Submit'
        }
    ];
    const formObj = new Form()
    const form = formObj.render({ inputs: inputs, formId: 'signupForm' })
    formblock.appendChild(form);

    const pLink = document.createElement('p');
    form.appendChild(pLink);

    const link = createLinkMenu()
    pLink.appendChild(link);

    form.addEventListener('submit', function (event) {

        event.preventDefault();

        const email = form.elements['email'].value;
        const password = form.elements['password'].value;
        const password_repeat = form.elements['password_repeat'].value;

        if (password !== password_repeat) {
            alert('Passwords is not equals');
            return;
        }
        if (email == "") {
            alert("Enter email!")
            return
        }

        httpReq.doPost({
            callback(res) {
                if (res.status == 208){
                    alert("Email already exist");
                    return;
                }
                if (res.status > 300) {
                    alert("Something was wrong");
                    return;
                }
                createProfile();
            },
            url: '/signup',
            data: {
                "email":email,
                "password":password,
            }
        });
    });
}

function createLeaders() {
    root.innerHTML = '';

    const header = document.createElement('div');
    header.dataset.sectionName = 'header';
    header.id = "header";
    root.appendChild(header);

    const logo = document.createElement('span');
    logo.id = 'logo';
    const headerTitle = document.createElement('h1');
    headerTitle.textContent = 'Leaderboard';
    headerTitle.id = "headerTitle";

    header.appendChild(logo);
    header.appendChild(headerTitle);

    const body = document.createElement('div');
    body.id = 'body';

    const pLink = document.createElement('p');
    headerTitle.appendChild(pLink);
    const link = createLinkMenu()
    pLink.appendChild(link);

    const table = document.createElement('table');

    table.border = 1;

    const tableHeader = document.createElement('tr');
    const th1 = document.createElement('th');
    const th2 = document.createElement('th');
    const th3 = document.createElement('th');

    th1.innerText = 'Nick';
    th2.innerText = 'Score';
    th3.innerText = 'Age';

    tableHeader.appendChild(th1);
    tableHeader.appendChild(th2);
    tableHeader.appendChild(th3);

    table.appendChild(tableHeader);

    /*  const em = document.createElement('em');
      em.textContent = 'Nothing to display';
      body.appendChild(em);*/

    xhr.open('POST', '/liderboards', true);
    xhr.setRequestHeader('Content-Type', 'text/plain; charset=utf-8');
    xhr.send('Request');

    xhr.onreadystatechange = function () {
        if (xhr.readyState !== 4) return;
        if (xhr.status !== 200) {
            alert(xhr.status + ': ' + xhr.statusText);
        } else {
            const top = JSON.parse(xhr.responseText);
            /*const d1 = document.createElement('div');
            const d = document.createElement('h1');

            d.innerText = text;
            d1.appendChild(d);

            root.appendChild(d1);*/

            const tbody = document.createElement('tbody');

            let username;
            let score;

            Object.entries(top).forEach(function ([id, info]) {
                username = info.nickname;
                score = info.score;


                const tr = document.createElement('tr');
                const tdUsername = document.createElement('td');
                const tdScore = document.createElement('td');
                // const tdAge = document.createElement('td');

                tdUsername.textContent = username;
                tdScore.textContent = score;
                //   tdAge.textContent = age;

                tr.appendChild(tdUsername);
                tr.appendChild(tdScore);
                // tr.appendChild(tdAge);

                tbody.appendChild(tr);

                table.appendChild(tbody);

            });
        }
    };

    body.appendChild(table);
    root.appendChild(body);

}

function createProfile() {
    root.innerHTML = '';
    const logo = document.createElement('div');
    logo.classList.add("p_name");

    const logo_header = document.createElement('h1');
    logo_header.innerHTML = 'Profile';
    logo.appendChild(logo_header);

    root.appendChild(logo);

    const body = document.createElement('div');
    body.id = 'body';
    root.appendChild(body);

    const pLink = document.createElement('p');
    body.appendChild(pLink);

    const link = createLinkMenu()
    pLink.appendChild(link);
}

function createAbout() {
    root.innerHTML = '';
    const header = document.createElement('div');
    header.id = "header";
    header.dataset.sectionName = 'header';
    root.appendChild(header);

    const logo = document.createElement('span');
    logo.id = 'logo';
    const headerTitle = document.createElement('h1');
    headerTitle.textContent = 'About';
    headerTitle.id = "headerTitle";

    header.appendChild(logo);
    header.appendChild(headerTitle);

    const body = document.createElement('div');
    body.id = 'body';

    const aboutTextblock = document.createElement('div');
    body.appendChild(aboutTextblock);

    const aboutText = document.createElement('p');
    aboutText.textContent = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
    aboutTextblock.appendChild(aboutText);

    const pLink = document.createElement('p');
    body.appendChild(pLink);

    const link = createLinkMenu()
    pLink.appendChild(link);


    root.appendChild(body);
}

const buttons = {
    signin: createSignIn,
    signup: createSignUp,
    leaders: createLeaders,
    profile: createProfile,
    about: createAbout,
    menu: createMenu,
};

root.addEventListener("click", function (event) {
    if (!(event.target instanceof HTMLAnchorElement)) return;

    event.preventDefault();

    const target = event.target;
    const href = target.dataset.href;

    root.innerHTML = '';
    console.log(href);
    buttons[href]();

});

createMenu();
