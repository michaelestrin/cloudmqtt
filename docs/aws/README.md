# michaelestrin/cloudmqtt/docs/aws -- Connecting to AWS IoT

I was able to connect the service to AWS IoT as follows.

First, in the AWS IoT console, I selected the __Register a thing__ option:

![Register A Thing](images/registerathing.png)

I then selected the __Create a single thing__ option:

![Select A Single Thing](images/selectasinglething.png)

In the subsequent screen, I then selected the __Create a single thing__ option:

![Create a Single Thing](images/createasinglething.png)

I then gave it a name and selected the __Next__ option:

![Name Your Thing](images/nameyourthing.png)

I then selected the __Create certification__ option for one-click certificate creation:

![Create Certificate](images/createcertificate.png)

Once the certificates were generated, I downloaded the certificate and keys.  I then activated the certificate by 
    selecting the __Activate__ option:
    
![Activate Certificate](images/activatecertificate.png)    

Once the certificate was activated, I selected the __Done__ option:

![Certificate Activated](images/certificateactivated.png)

I then created a new security policy by selecting the __Create__ option:

![Create New Policy](images/createnewpolicy.png)

I then named the policy, gave it _iot:*_ action permissions, checked the __Effect:Allow__ checkbox, and selected the 
    __Create__ option:
    
![Define the Policy](images/definethepolicy.png)

I then attached the policy to the certificate I created in an earlier step:

![Attach Policy](images/attachpolicy1.png)

![Attach Policy](images/attachpolicy2.png)

I then found my custom endpoint for connecting to the AWS IoT MQTT service:

![Custom Endpoint](images/customendpoint.png)

Finally, I modified my `configuration.toml` file to reference the certificate, the private key, and to incorporate 
    the custom endpoint:
    
![Configuration](images/configuration.png)