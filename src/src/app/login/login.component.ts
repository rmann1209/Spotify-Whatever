import { Component } from '@angular/core';
import { FormBuilder } from '@angular/forms'
import { User, loginUser } from './login.service';


@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  title = 'Boombox';

  constructor(private formBuilder:FormBuilder, private LoginUser:loginUser){}

  accountForm = this.formBuilder.group({
    Username:[''],
    Password:['']
  })

  loginUser(username: string, password: string) : void {
    this.LoginUser.loginUser({username, password} as User)
    .subscribe((response: any) => {
      console.log(response);
    });
  }

}
