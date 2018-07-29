package luascript

const (
	GET_BATCH_VALUES_LUA=`
         local key=''
		 local result={}
		 key=ARGV[1]
		 for i=1,#KEYS
		  do
		    result[i]=redis.call('hget',key,KEYS[i])
		  end
         return result
	`

	CHECK_GROUP_ALL_STATUS_LUA = `
		local result={}
		local status=''
		local argv=ARGV[1]
		local vccId=ARGV[2]
		local key=''
		local check=false
		result=redis.call('lrange',KEYS[1],0,-1)
		for i=1,#result
		do
		  key=string.format(argv,vccId,result[i])
		  status=redis.call('hget',key,'status')
		  if(status=='2')
		   then
			  check=true
			  break
		   end
		end
		return tostring(check)
	`
)
