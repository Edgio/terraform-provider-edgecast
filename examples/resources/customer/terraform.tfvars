# Provides values the variables used in main.tf

# Use the credentials provided to you by Verizon Media
credentials = {
    api_token = "AULdReDoB3gb0D7LNTx857NQvrcIKyvL"
    ids_client_secret = "CDbbMJw7FFJ11a7433ti1l9XgJHKr2Wk"
    ids_client_id = "31ef8e8f-0120-4112-8554-3eb11e83d58b"
    ids_scope = "ec.rules"
    api_address = "http://dev-api.edgecast.com"
    ids_address = "https://id-dev.vdms.io"
    api_address_legacy = "http://dev-api.edgecast.com"
}

new_customer_info = {
    company_name = "<Customer Name>"
    service_level_code = "STND"
    #all available services=> 1:HTTP Large Object,2:HTTPS Large Object,3:HTTP Small Object,4:HTTPS Small Object,6:Windows,7:Advanced Reports,8:Real-Time Stats,9:Token Auth,10:Edge Performance Analytics,15:Origin Storage,16:RSYNC,19:ADN,20:Download Manager,21:ADNS,22:Dedicated Hosting,23:Edge Optimizer,25:DNS Route,26:DNS Zones,29:DNS Health Checks,31:Bandwidth By Report Code,32:DNS-Standard,33:DNS-Adaptive,34:DNS-APR,38:WAF,39:Analysis Engine,40:HTTP Rate Limiting,41:Basic Rules v4.0,42:Advanced Rules v4.0,43:Mobile Device Detection Rules v4.0,44:Rules Engine v4.0,47:Translate,48:Dynamic Cloud Packaging,49:Encrypted HLS,50:Origin Shield,51:Reports and Logs,52:Log Delivery,54:SSA,56:Encrypted Key Rotation,57:Real-Time Log Delivery,58:Report Builder,59:Dynamic Imaging,60:China Delivery,61:WAF Essential,62:Report Builder Users,63:Report Builder Rows,64:Report Builder Reports,65:Edge Functions,66:Certificate Provisioning,67:Edge-Insights,68:Edge Image Optimizer,69:Url Redirects,70:Azure Cloud Storage
    services = [1,9,15,19]
    #available access modules => 1:Storage, 4:Analytics, 5:Admin, 7:Customer Origin, 8:Purge/Load, 21:Users, 22:Company, 25:Country Filtering, 26:Token Auth, 27:Dashboard, 29:HTTP Large, 30:Edge CNAMEs, 32:Core Reports, 40:Token Auth, 46:Token Auth, 53:Cache Settings, 56:HTTP Large Object, 71:HTTP Streaming, 72:ADN, 73:Customer Origin, 74:Purge/Load, 75:Token Auth, 76:Country Filtering, 77:Edge CNAMEs, 78:Cache Settings, 79:Application Delivery Network, 81:Tools, 138:Query-String Caching, 139:Query-String Logging, 140:Compression, 144:Query-String Caching, 145:Query-String Logging, 146:Compression, 149:Smooth Streaming Player, 153:JW Player, 157:Raw Log Settings, 159:Traffic Summary, 160:Bandwidth, 161:Data Transferred, 162:Hits, 163:Cache Statuses, 164:Cache Hit Ratio, 166:CDN Storage, 168:Notes, 169:HTTP Large, 170:HTTPS Large, 171:HTTP Small, 172:HTTPS Small, 174:Flash, 175:ADN, 176:ADN SSL, 177:HTTP Large, 178:HTTPS Large, 179:HTTP Small, 180:HTTPS Small, 182:Flash, 183:ADN, 184:ADN SSL, 185:All Platforms, 186:HTTP Large, 187:HTTP Small, 189:Flash, 190:ADN, 191:All Platforms, 192:HTTP Large, 193:HTTP Small, 194:ADN, 195:All Platforms, 196:HTTP Large, 197:HTTP Small, 198:ADN, 204:Usage, 386:IPv4/IPv6, 387:Data Transferred, 409:Custom Reports, 410:Edge CNAMEs, 411:Notes, 412:All Platforms, 413:HTTP Large, 414:HTTP Small, 415:Flash, 416:ADN, 479:Token Generator, 501:Add Users, 502:Edit Users
    access_modules = [1,4,5,7,8,21,22,25,26,27,29,30,32,40,46,53,56,71,72,73,74,75,76,77,78,79,81,138,139,140,144,145,146,149,153,157,159,160,161,162,163,164,166,168,169,170,171,172,174,175,176,177,178,179,180,182,183,184,185,186,187,189,190,191,192,193,194,195,196,197,198,204,386,387,409,410,411,412,413,414,415,416,479,501,502]
    delivery_region = 1 # TODO use data source to get ID
}
