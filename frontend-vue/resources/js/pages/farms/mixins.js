import stubFarm from '../../stores/stubs/farm'
import stubMessage from '../../stores/stubs/message'
import { mapActions, mapGetters } from 'vuex'

export default {
  data () {
    return {
      message: Object.assign({}, stubMessage),
      farm: Object.assign({}, stubFarm)
    }
  },

  computed: {
    ...mapGetters({
      countries: 'getCountries',
      cities: 'getCities',
      types: 'getAllFarmTypes'
    })
  },

  props: {
    intro: { type: Boolean, default: false}
  },

  methods: {
    ...mapActions([
      'createFarm',
      'fetchCitiesByCountryCode'
    ]),

    validateBeforeSubmit () {
      this.$validator.validateAll().then(result => {
        if (result) {
          this.create()
        }
      })
    },

    create () {
      // TODO if farm have UID we need to update endpoint API
      //
      this.createFarm(this.farm)
        .then(data => this.redirectTo(data))
        .catch(({ data }) => {
          this.message = data
        })
    },

    findMe () {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(position => {
          this.farm.latitude = position.coords.latitude
          this.farm.longitude = position.coords.longitude
        }, error => console.log(error))
      }
    },

    countryChanged (countryCode) {
      this.fetchCitiesByCountryCode(countryCode)
    }
  }
}
