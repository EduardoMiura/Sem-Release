import axios from 'axios';

const URL_REPOSITORY = 'http://localhost:9090/semrelease/repository';

const REPOSITORY_METHODS = {

    getRepositories: () => {
        const url = `${URL_REPOSITORY}`;
        return axios.get(url);
    },

};

export default REPOSITORY_METHODS;