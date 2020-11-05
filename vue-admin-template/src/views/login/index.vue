<template>
  <div class="login-container">
    <div class="login-box">

      <el-row>
        <el-col :span="12">
          
          <div class="login-aside-doc">
            <div class="login-aside-title">LiveRTC视频监控平台</div>
            <div class="login-aside-desc">LiveRTC是知名云计算及数据服务提供商，视频点播、互动直播及大规模异构数据的智能分析与处理等技术深度投入，致力以数据科技驱动数字化未来，赋能各行业全面进入数据时代。</div>
          </div>
          <div class="login-aside-img">
              <img src="https://sso.qiniu.com/asserts/login-aside.svg" width="326px" height="auto">
          </div>
        </el-col>
        <el-col :span="12"><div class="grid-content bg-purple-light login-main">

            <div class="login-title">连接数据，重塑价值</div>
            <div class="login-sub-title">欢迎来到LiveRTC，请登录！</div>

            <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" auto-complete="on" label-position="left">

              <div style="width:350px;">

              <el-form-item prop="username" >
                <span class="svg-container">
                  <svg-icon icon-class="user" />
                </span>    
                <el-input ref="username" v-model="loginForm.username" placeholder="用户名" name="username" type="text" auto-complete="on"></el-input>
              </el-form-item>

            <el-form-item prop="password">
              <span class="svg-container">
                <svg-icon icon-class="password" />
              </span>
              <el-input
                :key="passwordType"
                ref="password"
                v-model="loginForm.password"
                :type="passwordType"
                placeholder="Password"
                name="password"
                tabindex="2"
                auto-complete="on"
                @keyup.enter.native="handleLogin"
              />
              <span class="show-pwd" @click="showPwd">
                <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
              </span>
            </el-form-item>

            <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleLogin">Login</el-button>
            </div>
            <div class="tips">
              <span style="margin-right:20px;">username: admin</span>
              <span> password: any</span>
            </div>

          </el-form>   
              
        </div></el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { validUsername } from '@/utils/validate'
import '@/assets/css/login.css'

export default {
  name: 'Login',
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!validUsername(value)) {
        callback(new Error('Please enter the correct user name'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error('The password can not be less than 6 digits'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: 'admin',
        password: '111111'
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur', validator: validateUsername }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }]
      },
      loading: false,
      passwordType: 'password',
      redirect: undefined
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          this.$store.dispatch('user/login', this.loginForm).then(() => {
            this.$router.push({ path: this.redirect || '/' })
            this.loading = false
          }).catch(() => {
            this.loading = false
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    }
  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg:#283443;
$light_gray:#fff;
$cursor: #fff;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {

}
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $light_gray;
      height: 47px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
    }
  }


  .login-box{
    width:1024px;
    display: flex;
    -webkit-box-orient: vertical;
    -webkit-box-direction: normal;
    -ms-flex-direction: column;
    flex-direction: column;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    -webkit-box-pack: center;
    -ms-flex-pack: center;
    justify-content: center;

    border-radius: 4px 0px 0px 4px;
    background-color: rgba(0,170,231,1);
  }

  .bg-purple-left {
    width:100%;height:100%;
      
  }

  .login-aside-doc {
    height: 280px;
    padding-top: 84px;
  }
  .login-aside-doc .login-aside-title {
    text-align: center;
    line-height: 48px;
    font-size: 32px;
    font-weight: 400;
    color: rgba(255,255,255,1);
    margin-bottom: 12px;
  }

  .login-aside-doc .login-aside-desc {
    text-align: center;
    line-height: 18px;
    font-size: 12px;
    font-weight: 400;
    color: rgba(255,255,255,1);
    margin-left: 82px;
    margin-right: 81px;
  }

  .login-aside-img {
    width: 512px;
    height: 280px;
    padding: 1px 92px 87px 94px;
  }


  .login-form {
    padding-left: 0;
    padding-right: 0;
    border-radius: 0px 4px 4px 0px;
    box-shadow: 0px 20px 40px 0px rgba(0,0,0,0.08);
    margin-left: auto;
    margin-right: auto;
  }

  .login-main {
      background: rgba(255,255,255,1);
      height:600px;
      display: flex;
      -webkit-box-orient: vertical;
      -webkit-box-direction: normal;
      -ms-flex-direction: column;
      flex-direction: column;
      -webkit-box-align: center;
      -ms-flex-align: center;
      align-items: center;
      -webkit-box-pack: center;
      -ms-flex-pack: center;
      justify-content: center;
  }

  .login-main .login-title {
    color: rgba(0,170,231,1);
    font-size: 32px;
    font-weight: 400;
    line-height: 48px;
    margin-bottom: 12px;
    text-align: center;
  }
  .login-main .login-sub-title {
    color: rgba(0,0,0,0.45);
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
    margin-bottom: 64px;
    text-align: center;
  }

  .login-form .form-input {
    background: none;
    padding-left: 40px;
}

.form-control {
    display: block;
    width: 100%;
    height: 34px;
    padding: 6px 12px;
    font-size: 14px;
    line-height: 1.42857143;
    color: #555;
    background-color: #fff;
    background-image: none;
    border: 1px solid #ccc;
    border-radius: 4px;
    -webkit-box-shadow: inset 0 1px 1px rgba(0,0,0,.075);
    box-shadow: inset 0 1px 1px rgba(0,0,0,.075);
    -webkit-transition: border-color ease-in-out .15s,-webkit-box-shadow ease-in-out .15s;
    -o-transition: border-color ease-in-out .15s,box-shadow ease-in-out .15s;
    transition: border-color ease-in-out .15s,box-shadow ease-in-out .15s;
}

  .form {
    margin: auto;
    height: 560px;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }


  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    color: #454545;
  }
}
</style>

<style lang="scss" scoped>
$bg:#aaa;
$dark_gray:#889aa4;
$light_gray:#eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;
  display: flex;
  /*垂直排列*/
  flex-direction: column;
  align-items:center;/*由于flex-direction: column，因此align-items代表的是水平方向*/
  justify-content: center;/*由于flex-direction: column，因此justify-content代表的是垂直方向*/


  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
}
</style>
