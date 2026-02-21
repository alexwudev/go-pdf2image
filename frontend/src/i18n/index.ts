import { createI18n } from 'vue-i18n'
import zhTW from './zh-TW.json'
import en from './en.json'

export default createI18n({
  legacy: false,
  locale: 'zh-TW',
  fallbackLocale: 'en',
  messages: { 'zh-TW': zhTW, en },
})
